module register #(
    parameter N = 2
) (
    input  logic [N-1:0] data_in,
    input  logic         write, read,
    output logic [N-1:0] data_out,
);

    logic [N-1:0] data;

    always_ff @(posedge write) begin
        if (write)
            data <= data_in;
    end

    always_comb begin
        if (read)
            data_out <= data;
        else
            data_out <= 0;
    end

endmodule