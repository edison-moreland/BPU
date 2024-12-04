module register_rst #(
    parameter N = 3
) (
    input  logic         write, read, rst,
    input  logic [N-1:0] data_in,
    output logic [N-1:0] data_out,
);

    logic [N-1:0] data;

    always_ff @(posedge write, posedge reset) begin
        if (rst)
            data <= 0;
        else if (write)
            data <= data_in;
    end

    always_comb begin
        if (read)
            data_out <= data;
        else
            data_out <= 0;
    end

endmodule