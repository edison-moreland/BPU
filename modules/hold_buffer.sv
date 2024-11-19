module hold_buffer #(
    parameter N = 1
) (
    input  logic [N-1:0] data_in,
    input  logic         set_ready, set_data,
    output logic [N-1:0] data_out,
    output logic         ready,
);

    // reg ready;

    always_ff @(posedge set_ready)
        ready <= 1;

    always_ff @(posedge set_data, posedge ready) begin
        // if (set_ready) begin
        //     ready <= 1;
        // end

        if (ready && set_data) begin
            data_out <= data_in;
            ready <= 0;
        end
    end

    // always_comb
    //     is_ready <= ready;

endmodule
